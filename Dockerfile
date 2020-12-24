FROM chromedp/headless-shell:latest
ENV DEBIAN_FRONTEND=noninteractive
COPY ./ ./
RUN apt-get update && apt-get -y install tini golang git
RUN CGO_ENABLED=0 GOOS=linux go build -a
ENTRYPOINT ["tini", "--"]
CMD ["./pure-gym-tracker"]
