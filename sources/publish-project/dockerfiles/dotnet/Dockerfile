FROM mcr.microsoft.com/dotnet/core/aspnet:{{.Version}}
WORKDIR /app
COPY . .
ENTRYPOINT ["dotnet", "{{.Executable}}.dll"]