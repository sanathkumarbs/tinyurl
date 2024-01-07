FROM alpine

# Create user and set ownership and permissions as required
RUN mkdir /app
RUN adduser -D admin && chown -R admin /app

USER admin
COPY --chown=admin:admin bin/tiny /app
CMD ["/app/tiny"]