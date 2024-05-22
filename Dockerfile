FROM alpine as setup
RUN addgroup --gid 10000 -S appgroup && \
    adduser --uid 10000 -S appuser -G appgroup

FROM scratch as production
COPY --from=setup /etc/passwd /etc/passwd
COPY mirror /mirror
USER appuser

ENTRYPOINT ["/mirror"]
