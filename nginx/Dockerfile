FROM nginx:alpine

# Copy the frontend static files into the Nginx web root
COPY frontend/ /usr/share/nginx/html/

# Copy the Nginx config
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80 for the web gateway
EXPOSE 80

# Start Nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]
