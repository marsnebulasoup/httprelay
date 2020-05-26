/// <reference path="handler-response.ts" />
/// <reference path="handler-request.ts" />

//namespace HttpRelay.Proxy {
    type PlainHeaders = Headers | Record<string, string>
    type RouteParams = string[]

    interface RespondParams {
        body?: any,
        status?: number
        headers?: PlainHeaders
        fileName?: string
        download?: boolean
    }

    class HandlerCtx {
        constructor(
            public readonly request: HandlerRequest,
            public readonly abortSig: AbortSignal,
            public readonly routeParams: RouteParams
        ) {}

        get serverId(): string {
            return this.request.headerValue('HttpRelay-Proxy-ServerId')
        }

        get jobId(): string {
            return this.request.headerValue('HttpRelay-Proxy-JobId')
        }

        public respond(result: RespondParams = {}): HandlerResponse {
            return new HandlerResponse(result.body, result.status, result.headers, result.fileName, result.download)
        }
    }
//}