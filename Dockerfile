FROM scratch
COPY main ./
COPY zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
CMD ["./main"]