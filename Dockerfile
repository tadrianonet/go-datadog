# Etapa 1: Construção do binário Go
FROM golang AS builder

ENV     DD_ENV="env" \
        DD_VERSION="1.0.0" \
        DD_SERVICE="service"  \
        DD_PROFILING_ENABLED=true \
        DD_LOGS_INJECTION=true \
        DD_APPSEC_ENABLED=true \
        DD_IAST_ENABLED=true \
        DD_APPSEC_SCA_ENABLED=true

        # DD_AGENT_HOST="datadog-agent" \
        # DD_TRACE_AGENT_PORT="8126" 

# Definir o diretório de trabalho dentro do container
WORKDIR /app

# Copiar os arquivos de configuração do módulo Go
COPY go.mod go.sum ./

# Baixar as dependências necessárias
RUN go mod download

# Copiar o código-fonte do projeto
COPY . .

# Compilar a aplicação Go (gera o binário)
RUN GOOS=linux GOARCH=amd64 go build -o circuit-breaker main.go

# Etapa 2: Imagem final com o binário
FROM alpine:latest

# Definir variáveis de ambiente
ENV PORT=8080

# Definir o diretório de trabalho
WORKDIR /root/

# Copiar o binário da etapa de build
COPY --from=builder /app/circuit-breaker .

# Expor a porta da aplicação
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["./circuit-breaker"]
