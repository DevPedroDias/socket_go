# Use a imagem base do Golang
FROM golang:latest

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie o código-fonte para o diretório de trabalho
COPY . .

# Baixe as dependências do Go
RUN go mod download

# Compile o aplicativo
RUN go build -o main .

# Exponha a porta 5050 para a comunicação externa
EXPOSE 4545

# Execute o aplicativo quando o contêiner for iniciado
CMD ["./main"]