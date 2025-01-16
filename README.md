# Rate Limiter

O Rate Limiter é um middleware que limita o número de requisições que um cliente (por IP ou API_KEY) pode fazer em um determinado período de tempo. Ele é útil para proteger APIs contra abusos, como ataques de negação de serviço (DoS) ou uso excessivo de recursos.

## Como Funciona?

O Rate Limiter funciona com base em duas estratégias principais:

1. **Limite por IP**:
   - Restringe o número de requisições que um endereço IP pode fazer por segundo.
   - Se o limite for excedido, o IP é bloqueado temporariamente.

2. **Limite por API_KEY**:
   - Restringe o número de requisições que um token (API_KEY) pode fazer por segundo.
   - Se o limite for excedido, o token é bloqueado temporariamente.

### Estratégia de Bloqueio:
- Quando um cliente (IP ou API_KEY) excede o limite de requisições, ele é bloqueado por um tempo configurável (`BLOCK_TIME`).
- Após o tempo de bloqueio expirar, o cliente pode fazer novas requisições.

## Configuração

O Rate Limiter é configurado por meio de um arquivo `.env` e variáveis de ambiente. Abaixo estão as configurações disponíveis:

### Arquivo `.env`:
Crie um arquivo `.env` dentro da pasta cmd do projeto com as seguintes variáveis:

```env
REDIS_ADDR="redis:6379"         # Endereço do Redis
LIMIT_PER_IP=true               # Habilita o limite por IP (true/false)
LIMIT_PER_API_KEY=false         # Habilita o limite por API_KEY (true/false)
MAX_REQUESTS_PER_SECOND=10      # Número máximo de requisições por segundo, deve ser um valor inteiro
BLOCK_TIME=2                    # Tempo de bloqueio em segundos, deve ser um valor inteiro
HTTP_CODE_LIMIT_REACHED=429     # Código HTTP retornado quando o limite é excedido
MESSAGE_LIMIT_REACHED="Too many requests"  # Mensagem retornada quando o limite é excedido
```

### Explicação das Variáveis:
- `REDIS_ADDR`: Endereço do Redis (usado para armazenar contadores e bloqueios).
- `LIMIT_PER_IP`: Habilita ou desabilita o limite por IP.
- `LIMIT_PER_API_KEY`: Habilita ou desabilita o limite por API_KEY.
- `MAX_REQUESTS_PER_SECOND`: Número máximo de requisições permitidas por segundo.
- `BLOCK_TIME`: Tempo (em segundos) que um cliente fica bloqueado após exceder o limite.
- `HTTP_CODE_LIMIT_REACHED`: Código HTTP retornado quando o limite é excedido (padrão: 429 - Too Many Requests).
- `MESSAGE_LIMIT_REACHED`: Mensagem retornada quando o limite é excedido.


## Como Rodar a Aplicação
### Pré-requisitos
- Docker e Docker Compose instalados.
- Arquivo `.env` configurado (veja a seção de configuração acima).

### Passos para Rodar a Aplicação
1. **Clone o Respositório**
```bash
git clone https://github.com/Ruteski/rate-limiter.git
cd rate-limiter
```

2. **Construa e Suba os Contêineres:**
Execute o seguinte comando para construir a imagem da aplicação e subir os contêineres (Redis e aplicação):
```bash
docker-compose up --build -d
```

3. **Verifique os Contêineres**
Verifique se os contêineres estão rodando:
```bash
docker ps
```
Você deve ver dois contêineres:
- `redis_rate_limiter`: Contêiner do Redis.
- `rate_limiter_app`: Contêiner da aplicação.

4. **Teste a Aplicação**
Use o comando curl para testar a aplicação:
```bash
for i in {1..12}; do curl -H "api_key:qwe-987" -i localhost:8080; done
```
- As primeiras 10 requisições serão permitidas (configurado em `MAX_REQUESTS_PER_SECOND`).
- As requisições excedentes retornarão o código HTTP `429` e a mensagem `"Too many requests"` (configurado em `MESSAGE_LIMIT_REACHED`)

5. **Verifique os Logs**
Para ver os logs da aplicação, use:
```bash
docker logs rate_limiter_app
```

6. **Pare os Contêineres:**
Para parar os contêineres, use:
```bash
docker-compose down
```

## Estrutura do Projeto
A estrutura do projeto é a seguinte:
```
/projeto
  /cmd
    main.go          # Ponto de entrada da aplicação
    .env             # Arquivo de configuração
  /config
    config.go        # Código para carregar as configurações
  /limiter
    limiter.go       # Lógica do Rate Limiter
  /middleware
    middleware.go    # Middleware para integrar o Rate Limiter
  /store
    in-memory.go     # Implementação do armazenamento em memória
    redis.go         # Implementação do armazenamento no Redis
    store.go         # Interface do armazenamento
  /tests
    rate_limiter_test.go  # Testes automatizados
  docker-compose.yaml     # Configuração do Docker Compose
  Dockerfile              # Dockerfile para construir a imagem da aplicação
  go.mod                  # Dependências do projeto
  go.sum                  # Checksum das dependências
```

## Teste Automatizados
O projeto inclui testes automatizados para garantir que o Rate Limiter funcione corretamente. Para rodar os testes, execute na raiz do projeto:
```bash
go test -v ./tests
```
- Os testes cobrem cenários como:
- Limite de requisições por IP.
- Limite de requisições por API_KEY.
- Bloqueio temporário após exceder o limite.
- Liberação após o tempo de bloqueio expirar.

## Como Contribuir
1. Faça um fork do repositório.
2. Crie uma branch para sua feature ou correção:
```bash
git checkout -b minha-feature
```
3. Envie suas alterações:
```bash
git commit -m "Adiciona nova feature"
git push origin minha-feature
```
4. Abra um pull request no repositório original.


## Conclusão
O Rate Limiter é uma solução eficiente para proteger APIs contra abusos e garantir o uso justo dos recursos. Com configurações flexíveis e suporte a Redis ou armazenamento em memória, ele pode ser adaptado para diferentes cenários de uso.

Para mais informações, consulte o código-fonte ou entre em contato com os mantenedores do projeto. 🚀