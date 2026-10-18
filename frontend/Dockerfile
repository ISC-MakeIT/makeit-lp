FROM node:22-alpine
WORKDIR /app
# node_modules は名前付きボリューム。初回/依存変更時は起動時にinstallする
CMD ["sh", "-c", "npm ci && npm run dev"]
