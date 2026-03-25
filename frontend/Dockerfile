# Stage 1: Build the React app
FROM node:20-alpine AS build
WORKDIR /app

# Install dependencies first (better caching)
COPY package*.json ./
RUN npm install

# Copy the rest of the code and build
COPY . .

ARG Backend_API_URL
ENV Backend_API_URL=$Backend_API_URL

RUN npm run build

# Stage 2: Serve with Nginx
FROM nginx:stable-alpine

# Copy the build output from Stage 1 to Nginx's public folder
COPY --from=build /app/dist /usr/share/nginx/html

# Expose port 80 (standard for HTTP)
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]



# ```bash
# docker build \
# --build-arg VITE_Backend_API_URL=http://10.42.206.173:30008 \
# -t surajmane121045/task-frontend:1.1 .
# ```
#  docker build -t my-frontend .
# docker run -d --name my-frontend --network my-net -p 82:80 my-frontend
