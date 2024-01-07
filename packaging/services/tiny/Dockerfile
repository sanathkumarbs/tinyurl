FROM alpine as builder
# Create user and set ownership and permissions as required
RUN mkdir /app
RUN adduser -D admin && chown -R admin /app
COPY --chown=admin:admin bin/tiny /app

FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app /
USER admin
CMD ["/tiny"]