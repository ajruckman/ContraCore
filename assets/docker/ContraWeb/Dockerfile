FROM mcr.microsoft.com/dotnet/core/aspnet:3.1.3-buster-slim

#FROM ubuntu:eoan

#RUN apt update
#RUN apt install apt-transport-https ca-certificates

#COPY packages-microsoft-prod.deb /opt/

#RUN dpkg -i /opt/packages-microsoft-prod.deb


#COPY aspnetcore-runtime-3.1.3-linux-musl-x64.tar.gz /opt/
COPY publish/* /opt/contraweb/

WORKDIR /opt/contraweb/

EXPOSE 80/tcp

#CMD /bin/sh
CMD ["dotnet", "./Web.dll"]
