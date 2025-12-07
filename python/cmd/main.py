import os
import uvicorn
from python.routes.routes import setup_app
from python.utils.db import init_db_from_env

def main():
    app_port = os.getenv("APP_PORT", ":8992")
    db = init_db_from_env()
    app = setup_app(db)
    host = "0.0.0.0"
    port = int(app_port.strip().lstrip(":"))
    uvicorn.run(app, host=host, port=port)

if __name__ == "__main__":
    main()
