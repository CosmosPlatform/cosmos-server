# cosmos-server
Server for the cosmos platform

## File Structure

```
cosmos-server
├── pkg
│   ├── app                  // Package where the app is defined and has the functions to be initialized.
│   ├── api                  // Package where the API is defined
│   │   ├── routes           // Package where the routes are defined.
│   │   └── dto              // Package where the data transfer objects are defined.
│   ├── config               // Package where the configuration format is defined.
│   ├── server               // Package where the server is defined.
│   ├── services             // Package where the services are defined.
│   └── test                 // Package where test utilities are defined
└── config                   // Configuration files
```