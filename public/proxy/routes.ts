/// <reference path="route.ts" />
/// <reference path="handler-ctx.ts" />

//namespace HttpRelay.Proxy {

    interface SelectedRoute {
        handlerFunc: HandlerFunc,
        params: RouteParams
    }

    class Routes {
        private readonly routes: Route[] = []
        public notFoundHandlerFunc: HandlerFunc

        constructor() {
            this.notFoundHandlerFunc = ctx => ctx.respond({
                status: 404,
                body: `Handler not found for the "${ctx.request.method} ${ctx.request.path}" route on "${ctx.serverId}" server.`
            })
        }

        public add(method: string, path: string, handlerFunc: HandlerFunc): void {
            let route = new Route(method, path, handlerFunc)
            this.routes.push(route)
            this.routes.sort((a, b) => a.compare(b))
        }

        public find(method: string, path: string): SelectedRoute {
            let routeParams: RouteParams = []
            let route = this.routes.find(r => {
                let matchRes = r.match(method, path)
                if (matchRes != null) {
                    routeParams = matchRes
                    return true
                }
            })

            return <SelectedRoute> {
                handlerFunc: route ? route.handlerFunc : this.notFoundHandlerFunc,
                params: routeParams
            }
        }

        public get(path: string, handlerFunc: HandlerFunc): void {
            this.add('GET', path, handlerFunc)
        }

        public post(path: string, handlerFunc: HandlerFunc): void {
            this.add('POST', path, handlerFunc)
        }
    }
//}