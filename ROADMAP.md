# Marble Roadmap

## Features

- Hash Literal

- Closures bound to object instances.

  This will lay the groundwork for allowing closures to be named & bound to classes. As a Class will be an instances of an object.

  The bind method will take 3 arguments:

    1) The object instance the named closure will be bound to

    2) An identifier for the closure

    3) The closure to be called when the calling the closure from the object instance. The closure MUST take at least 1 arguments, a `self` that will be the instance of the bound object.

  Example:

  ```marble
  const message = "Foo Bar"

  bind(message, length, fn(self) { len(self) })

  message.length

  // => Int(7)
  ```

- Classes

- Namespaces

- API Binding

  API binding will allow a host application to provide native code bindings and functions to the evaluation runtime.
