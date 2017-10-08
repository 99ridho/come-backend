# REST API Reference

## Login

* Auth : No
* Request (`Content-Type: application/json`)
    
    ```
    {
        "email" : "your_email",
        "password: "your_password"
    }
    ```
* Response (`Content-Type: application/json`)
    
    ```
    {
        "status": "status", // depend to login result, "success" or "failed"
        "message": "message", // depend to login result
        "token": "token" // JWT token will present if logged in successfully
    }
    ```