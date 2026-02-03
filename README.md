# OpenCLI

CLI para interagir com o OpenProject.

## Instalação

```bash
go build -o op .
```

## Configuração

Crie o diretório e o arquivo de configuração:

```bash
mkdir -p ~/.config/opcli
```

Depois crie o arquivo `~/.config/opcli/config.yaml`:

```yaml
base_url: https://seu-openproject.com
api_key: sua-api-key-aqui
project: nome-do-projeto
```

### Obtendo a API Key

1. Acesse seu OpenProject
2. Vá em **My Account** > **Access tokens**
3. Crie um novo token de API
4. Copie o token gerado e coloque no `api_key`

### Variáveis de ambiente (alternativa)

Você também pode configurar via variáveis de ambiente:

```bash
export OPENPROJECT_BASE_URL=https://seu-openproject.com
export OPENPROJECT_API_KEY=sua-api-key-aqui
export OPENPROJECT_PROJECT=nome-do-projeto
```

## Comandos

### Listar Work Packages

```bash
op wp list              # lista com paginação (70 por página)
op wp list --all        # lista todos
op wp list --page 2     # página específica
op wp list --size 20    # define itens por página
```

### Ver detalhes de um Work Package

```bash
op wp show 123
```

### Atribuir Work Package a você

```bash
op wp assign-me 123
```

## Roadmap

Ideias em desenvolvimento:

- Integração com Ollama para criar tasks a partir de screenshots/prints
- Criação de Work Packages via IA com descrição automática
- Mais comandos para gerenciamento de projetos

## Contribuindo

Toda ideia de melhoria é bem-vinda! Sinta-se à vontade para abrir issues ou enviar PRs.

## Autor

Desenvolvido por [Guilherme Alves](https://www.guialves.site)
