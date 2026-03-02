FROM node:22-alpine

ENV NODE_ENV development

# npm is built into the node image, so yarn is removed.
# python2 is updated to python3, which is required for modern node-gyp builds in Node 22.
RUN apk --no-cache add python3 make gcc g++
ADD start.sh /
RUN chmod +x /start.sh

CMD ["/start.sh"]