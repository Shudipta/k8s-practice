# our base image
FROM ubuntu

# copy the build file
#RUN ["mkdir", "/book_server"]
COPY . /samp-crd-ctl
WORKDIR /samp-crd-ctl

# specify the port number the container should expose
EXPOSE 10000

# run the server
# CMD ["cd", "/samp-crd-ctl"]
CMD ["./main", "create"]
