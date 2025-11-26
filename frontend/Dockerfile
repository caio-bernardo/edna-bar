# Multi-stage Dockerfile for building a JS frontend and serving it with nginx
# - Build stage uses Node to install deps and produce a production build
# - Runner stage uses nginx to serve the static files with an SPA-friendly config
#
# Usage:
#  docker build -t edna-frontend:latest .
#  docker run --rm -p 80:80 edna-frontend:latest
#
# Notes:
# - Expects a standard `npm run build` script that outputs to `dist/` (Vite and many other setups).
# - If your project builds to a different directory (e.g. `build/`), update the COPY from the builder accordingly.

################################
# Builder stage
################################
FROM node:20-alpine AS builder

# set working dir
WORKDIR /app

# Bring PATH update for locally installed binaries
ENV PATH /app/node_modules/.bin:$PATH


# Copy package manifests first to take advantage of Docker layer caching
# Copy both package*.json and yarn.lock when present
COPY package.json package-lock.json* yarn.lock* ./

# If native builds are needed, keep build tools available in this stage
RUN apk add --no-cache python3 make g++ || true

# Install dependencies: try `npm ci` first (if lockfile present), fallback to `npm install`
RUN npm ci --silent --no-audit --prefer-offline || npm install --silent --no-audit

# Ensure NODE_ENV is set only after dependencies are installed so devDependencies (like Vite) are installed during build
ENV NODE_ENV=${NODE_ENV}

# Copy source
COPY . .

# Build the production assets. Use `--if-present` so this won't fail for projects without a build script.
RUN npm run build --if-present

################################
# Runner stage (Used for deploy)
################################
FROM nginx:stable-alpine AS runner

# Copy built assets from the builder stage to nginx html directory.
# Vite outputs to `/app/dist` by default — adjust if your build produces a different folder.
COPY --from=builder /app/dist /usr/share/nginx/html

# Use the default nginx foreground cmd for Docker images
CMD ["nginx", "-g", "daemon off;"]
