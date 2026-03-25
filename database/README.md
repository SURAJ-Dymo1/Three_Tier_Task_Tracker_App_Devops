Since you want a **custom Docker image for MongoDB**, the Dockerfile should start from the official MongoDB base image and optionally add initialization scripts.

Below is a **correct and clean custom Dockerfile**.

---

# 1️⃣ Dockerfile (Custom MongoDB Image)

```dockerfile
# Use official MongoDB image
FROM mongo:7

# Set default admin credentials
ENV MONGO_INITDB_ROOT_USERNAME=admin
ENV MONGO_INITDB_ROOT_PASSWORD=admin123
ENV MONGO_INITDB_DATABASE=mydatabase

# Copy initialization scripts (optional)
COPY init-mongo.js /docker-entrypoint-initdb.d/

# Expose MongoDB port
EXPOSE 27017
```

---

# 2️⃣ Initialization Script (Create User Automatically)

Create file **`init-mongo.js`**

```javascript
db = db.getSiblingDB("mydatabase");

db.createUser({
  user: "appuser",
  pwd: "password123",
  roles: [
    {
      role: "readWrite",
      db: "mydatabase"
    }
  ]
});

db.createCollection("sample_collection");
```

This will automatically run when the container starts.

---

# 3️⃣ Build Your Custom Image

Inside the folder containing the Dockerfile:

```bash
docker build -t my-mongo .
```

---

# 4️⃣ Run MongoDB Container

```bash
docker run -d \
--name my-mongo \
--network my-net \
-p 27017:27017 \
-v mongo_data:/data/db \
my-mongo
```

---

# 5️⃣ Verify MongoDB Running

```bash
docker ps
```

You should see:

```
my-mongo
0.0.0.0:27017->27017/tcp
```

---

# 6️⃣ Connect to MongoDB

Run:

```bash
docker exec -it mongodb mongosh -u admin -p admin123 --authenticationDatabase admin
```

---

# 7️⃣ Test Database

Inside the shell:

```javascript
use mydatabase
show collections
```

You should see:

```
sample_collection
```

---

# 8️⃣ Connect From Your Backend

Example connection string:

```
mongodb://appuser:password123@localhost:27017/mydatabase
```

---

# 9️⃣ Persist Data (Important)

You already mounted volume:

```
mongo_data:/data/db
```

Check volumes:

```bash
docker volume ls
```

---

# 🔟 Folder Structure (Recommended)

```
project/
│
├── Dockerfile
├── init-mongo.js
```

---

# 🚀 Pro DevOps Tip

For projects like your **backend + ML pipeline**, it's better to run MongoDB with **Docker Compose** so everything runs together:

```
Backend (Flask)
MongoDB
Mongo Express
```

---

✅ If you want, I can also give you a **complete DevOps project structure**:

```
frontend/
backend/
mongodb/
docker-compose.yml
nginx/
```

This is the **exact structure used in production full-stack deployments**.












































Here are the most common and reliable ways to check if **MongoDB** is running on Linux (as of 2026, most distributions use systemd):

### 1. Best & Recommended Method (modern Linux: Ubuntu 18.04+, Debian 10+, CentOS 8+/RHEL 8+, Fedora, etc.)

```bash
sudo systemctl status mongod
```

- What you'll see if it's **running**:
  ```
  ● mongod.service - MongoDB Database Server
     Loaded: loaded (/lib/systemd/system/mongod.service; enabled; vendor preset: enabled)
     Active: active (running) since ...
  ```

- If it's **not running**:
  ```
  ○ mongod.service - MongoDB Database Server
     Loaded: loaded (...)
     Active: inactive (dead)
  ```
  or
  ```
  × mongod.service - MongoDB Database Server
     Active: failed (...)
  ```

Quick one-liner check (returns 0 = running, non-zero = not running):

```bash
sudo systemctl is-active mongod >/dev/null && echo "MongoDB is running" || echo "MongoDB is NOT running"
```

### 2. Alternative names people sometimes use

Sometimes the service is called `mongodb` instead of `mongod` (especially older Ubuntu/Debian packages):

```bash
sudo systemctl status mongodb
# or
sudo service mongodb status
```

Try both if the first one says "Unit mongod.service could not be found."

### 3. Very reliable process-based check (works even without systemd / when started manually)

```bash
# Method A - cleanest
pgrep -f mongod

# Method B - more detailed
ps aux | grep -i '[m]ongod'

# Method C - count running instances
ps -ef | grep -i '[m]ongod' | wc -l
```

- If you see a number > 0 (or a line with `/usr/bin/mongod` or similar) → MongoDB is running
- If nothing → not running

### 4. Quick port check (default port = 27017)

```bash
# Using ss (recommended, faster than netstat)
ss -tuln | grep 27017

# or old-school netstat (if installed)
netstat -tuln | grep 27017
```

If you see something like `tcp LISTEN 0 128 0.0.0.0:27017` → MongoDB is listening (almost certainly running).

### 5. Try to connect (the ultimate test)

```bash
mongosh          # or mongo (older versions)
# or with explicit host
mongosh --host localhost --port 27017
```

If it connects and shows a prompt → definitely running.

If it says "connection refused" or times out → not running (or wrong port / bind IP / firewall).

### Summary Table – Pick One

| Method                            | Command                              | When to use                              | Most reliable? |
|-----------------------------------|--------------------------------------|------------------------------------------|----------------|
| Modern systemd                    | `sudo systemctl status mongod`       | Ubuntu 18+, Debian 10+, RHEL 8+, Fedora  | ★★★★★          |
| Process check                     | `pgrep mongod` or `ps aux │ grep [m]ongod` | Any Linux, manual start, containers      | ★★★★☆          |
| Service (older style)             | `sudo service mongod status`         | Older systems or SysV init               | ★★★☆☆          |
| Port listening                    | `ss -tuln │ grep 27017`             | Quick, no sudo needed                    | ★★★★☆          |
| Try to connect                    | `mongosh`                            | You want to be 100% sure it's usable     | ★★★★★          |

Most of the time just run:

```bash
sudo systemctl status mongod
```

That's usually all you need.

Good luck!


**Here’s the complete, clean way to set up MongoDB in Docker with a custom database created automatically on first start.**

We will use:
- Official `mongo:8` image (latest stable major version as of 2026)
- A small init script to **create your database + a ready-to-use user**
- A simple Dockerfile (as you asked)
- Named volume for data persistence

### Step 1: Create project folder
```bash
mkdir mongodb-setup && cd mongodb-setup
```

### Step 2: Create the database initialization script
Create a file named `mongo-init.js`:

```bash
cat > mongo-init.js << 'EOF'
print("=== MongoDB Initialization Started ===");

const dbName = "myappdb";          // ← CHANGE THIS TO YOUR DESIRED DATABASE NAME
const db = db.getSiblingDB(dbName);

db.createUser({
  user: "appuser",                 // ← Your app username
  pwd:  "appsecret123",            // ← Change this password!
  roles: [{ role: "readWrite", db: dbName }]
});

print(`✅ Database '${dbName}' created successfully!`);
print(`✅ User 'appuser' with readWrite access created!`);
EOF
```

### Step 3: Create the Dockerfile
Create a file named `Dockerfile`:

```dockerfile
FROM mongo:8

# Copy our init script so it runs automatically on first start
COPY mongo-init.js /docker-entrypoint-initdb.d/init.js

# (Optional) Expose the port for documentation
EXPOSE 27017
```

### Step 4: Build your custom MongoDB image
```bash
docker build -t my-mongodb .
```

### Step 5: Run the MongoDB container
```bash
docker run -d \
  --name my-mongodb \
  --network my-net \
  -p 27017:27017 \
  -e MONGO_INITDB_ROOT_USERNAME=admin \
  -e MONGO_INITDB_ROOT_PASSWORD=adminsecret123 \
  -e MONGO_INITDB_DATABASE=myappdb \
  -v mongodb-data:/data/db \
  -v $(pwd)/mongo-init.js:/docker-entrypoint-initdb.d/mongo-init.js:ro \
  my-mongodb
```

**What this does:**
- Creates root admin user (`admin` / `adminsecret123`)
- Runs your `mongo-init.js` **only on first start** → creates database `myappdb` + user `appuser`
- Persists all data in a Docker volume (`mongodb-data`)
- Maps port 27017 so you can connect from your host

### Step 6: Verify everything is working

```bash
# 1. Check container is running
docker ps | grep my-mongodb

# 2. See initialization logs (you should see your success messages)
docker logs my-mongodb | tail -30

# 3. Connect to MongoDB (as root admin)
docker exec -it my-mongodb mongosh -u admin -p adminsecret123

# Inside mongosh, run these commands:
show dbs;                    # ← You should see "myappdb" and "admin"
use myappdb;
db.getUsers();               # ← You should see "appuser"
```

### How to connect from your host/app (examples)

**Using mongosh (recommended):**
```bash
# As root admin
mongosh "mongodb://admin:adminsecret123@localhost:27017/?authSource=admin"

# As app user (recommended for your application)
mongosh "mongodb://appuser:appsecret123@localhost:27017/myappdb?authSource=myappdb"
```

**In your Node.js/Python/etc app, use this connection string:**
```
mongodb://appuser:appsecret123@localhost:27017/myappdb?authSource=myappdb
```

### Useful commands later

```bash
# Stop container
docker stop my-mongodb

# Start again (data is still there)
docker start my-mongodb

# Remove container (data stays in volume)
docker rm my-mongodb

# Delete everything and start fresh (re-init will run again)
docker rm -f my-mongodb && docker volume rm mongodb-data
```

**Done!**  
You now have a fully running MongoDB with your custom database `myappdb` created automatically.

Want to change the database name, username, or password? Just edit `mongo-init.js`, rebuild the image (`docker build -t my-mongodb .`), and run a fresh container.

Let me know if you want the same setup with **docker-compose.yml** (even easier) or want to add replica set / authentication tweaks!


To allow your backend service (running in the same custom network called `my-net`) to connect to MongoDB **using the container name** `my-mongodb` (DNS name resolution), you need to attach the MongoDB container to that user-defined network.

Add this option to your `docker run` command:

```bash
--network my-net
```

### Modified command (recommended)

```bash
docker run -d \
  --name my-mongodb \
  --network my-net \
  -p 27017:27017 \
  -e MONGO_INITDB_ROOT_USERNAME=admin \
  -e MONGO_INITDB_ROOT_PASSWORD=adminsecret123 \
  -v mongodb-data:/data/db \
  mongo
```

*(I also replaced `my-mongodb` at the end with the official image name `mongo` — assuming that's what you intended. If you're using a different tag like `mongo:7`, write that instead.)*

### Quick checklist – make sure these are true

1. The network already exists  
   ```bash
   docker network create my-net
   ```
   (run this **once** — you don't need to recreate it every time)

2. Your backend container is also started with  
   ```bash
   --network my-net
   ```

3. In your backend code, connect to MongoDB using:  
   ```
   mongodb://admin:adminsecret123@my-mongodb:27017
   ```
   or (with database name)
   ```
   mongodb://admin:adminsecret123@my-mongodb:27017/myappdb?authSource=admin
   ```

### Summary – what to add / change

| What to do                          | Flag to add/change             |
|:------------------------------------|:-------------------------------|
| Connect container to "my-net"       | `--network my-net`             |
| (optional but cleaner) image name   | `mongo` instead of `my-mongodb` at the end |

That's all you need. After this change, containers on `my-net` can reach each other by container name (like normal DNS in a small private network).

Good luck! 🚀