FROM node
ADD app.js package.json /app/
WORKDIR /app
RUN npm install
ENTRYPOINT [ "node" ]
CMD ["/app/app.js"]