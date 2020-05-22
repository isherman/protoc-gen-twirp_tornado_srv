# Twirp Tornado Server Generator

Generate [Tornado](https://www.tornadoweb.org/) `RequestHandler`s for
[Twirp](https://github.com/twitchtv/twirp) requests.

Based on [https://github.com/daroot/protoc-gen-twirp_python_srv](https://github.com/daroot/protoc-gen-twirp_python_srv).

## Install

```bash
go get -u github.com/isherman/protoc-gen-twirp_tornado_srv
```

## Requirements

Generated code requires:
- [tornado](https://pypi.org/project/tornado/)
- [protobuf](https://pypi.org/project/protobuf/)
- `validate/validator.py` from [protoc-gen-validate](https://github.com/envoyproxy/protoc-gen-validate)

`validate/validator.py` is not available on PyPI. It should be copied or symlinked into your `PYTHONPATH`.


## Usage

Generate messages and service stubs from `.proto` files

```bash
protoc --proto_path=protos --python_out=python/genproto --twirp_tornado_srv_out=python/gensrv protos/*_service.proto
```

Subclass `ServiceImpl` classes to provide service implementations

```python
from gensrv.acme.foo_service_twirp_srv import FooServiceImpl


class FooService(FooServiceImpl):
    def CreateFoo(self, create_foo_request):
        pass

    def ListFoos(self, list_foos_request):
        pass

    def GetFoo(self, get_foo_request):
        pass

    def DeleteFoo(self, delete_foo_request):
        pass
```

Register the generated request handlers and their corresponding service implementations with your Tornado application:

```python
# Service implementation
from acme.foo_service import FooService

# Generated service stub
from gensrv.acme.foo_service_twirp_srv import FooServiceRequestHandler

class Application(tornado.web.Application):
    def __init__(self):
        handlers = []
        for service, handler in [(FooService, FooServiceRequestHandler)]:
                handlers.append((rf'{handler.prefix}/.*/?', handler, {"service": service()}))

        super().__init__(handlers)

if __name__ == "__main__":
    app = Application()
    app.listen(8888)
    tornado.ioloop.IOLoop.current().start()
```

## TODO

- Provide protoc plugin options for
  - CORS
  - Validation
