package model

import (
	"time"
)

type McastSeq struct {
	mcastDataMap map[int]*McastData
	newSeqId     int
	oldSeqId     int
	comm
}

func NewMcastSeq(initialSeqId int) *McastSeq {
	initialMap := make(map[int]*McastData)
	initialMap[initialSeqId] = NewMcastData()
	return &McastSeq{
		mcastDataMap: initialMap,
		newSeqId:     initialSeqId,
		oldSeqId:     initialSeqId,
		comm:         newComm(),
	}
}

func (ms *McastSeq) Close() {
	ms.RLock()
	defer ms.RUnlock()
	for _, v := range ms.mcastDataMap {
		v.Close()
	}
}

func (ms *McastSeq) GetData(seqId int) (data *TeeData, ok bool) {
	ms.RLock()
	defer ms.RUnlock()

	if seqId < ms.newSeqId {
		if mcastData, ok := ms.mcastDataMap[seqId]; ok {
			return mcastData.data, ok
		}
	}

	return
}

func (ms *McastSeq) Read(wantedSeqId int, closeChan <-chan struct{}) (data *TeeData, seqId int, ok bool) {
	ms.AddWaiter()
	defer ms.RemoveWaiter()

	ms.Lock()
	if wantedSeqId == -1 {
		if ms.newSeqId == ms.oldSeqId {
			seqId = ms.newSeqId
		} else {
			seqId = ms.newSeqId - 1
		}
	} else if wantedSeqId > ms.newSeqId {
		seqId = ms.newSeqId
	} else if wantedSeqId < ms.oldSeqId {
		seqId = ms.oldSeqId
	} else {
		seqId = wantedSeqId
	}

	mcastData := ms.mcastDataMap[seqId]
	ms.Unlock()

	data, ok = mcastData.Read(closeChan)
	return
}

func (ms *McastSeq) Write(data *TeeData) (seqId int) {
	ms.Lock()
	defer ms.Unlock()

	ms.accessed = time.Now()
	seqId = ms.newSeqId
	mcastData := ms.mcastDataMap[seqId]
	mcastData.Write(data)
	ms.preserveSize()
	ms.newSeqId += 1
	ms.mcastDataMap[ms.newSeqId] = NewMcastData()
	return
}

func (ms *McastSeq) preserveSize() {
	for ms.size() > 11000000 { // Total allowed 11Mb while reqest limited 10Mb
		ms.removeOldest()
	}
}

func (ms *McastSeq) removeOldest() {
	if ms.oldSeqId < ms.newSeqId {
		ms.remove(ms.oldSeqId)
		ms.oldSeqId += 1
	}
}

func (ms *McastSeq) remove(seqId int) {
	if mcastData, ok := ms.mcastDataMap[seqId]; ok {
		delete(ms.mcastDataMap, seqId)
		mcastData.Close()
	}
}

func (ms *McastSeq) size() (size int) {
	for _, v := range ms.mcastDataMap {
		size += v.Size()
	}
	return
}

func (ms *McastSeq) Size() int {
	ms.RLock()
	defer ms.RUnlock()
	return ms.size()
}

func (ms *McastSeq) DataCount() int {
	ms.RLock()
	defer ms.RUnlock()
	return len(ms.mcastDataMap)
}

func (ms *McastSeq) NewSeqId() int {
	return ms.newSeqId
}
