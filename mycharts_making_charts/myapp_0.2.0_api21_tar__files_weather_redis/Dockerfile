FROM node
ADD app.js package.json /app/
ADD config app/config
WORKDIR /app
RUN npm install
ENTRYPOINT [ "node" ]
CMD ["/app/app.js"]
