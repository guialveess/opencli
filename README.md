# OpenCLI

CLI para interagir com o OpenProject.

## Stack

- **Go** - Linguagem principal
- **Cobra** - Framework para CLI
- **Viper** - Gerenciamento de configuração
- **Lipgloss** - Estilização do terminal
- **Ollama** - IA local para análise de imagens

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

### `op wp list`

Lista os Work Packages do projeto.

```bash
op wp list              # lista com paginação (70 por página)
op wp list --all        # lista todos
op wp list --page 2     # página específica
op wp list --size 20    # define itens por página
```

| Flag | Alias | Descrição |
|------|-------|-----------|
| `--all` | `-a` | lista todos os work packages |
| `--page` | `-p` | número da página |
| `--size` | `-s` | itens por página |

### `op wp show`

Exibe detalhes de um Work Package.

```bash
op wp show 123
```

### `op wp assign-me`

Atribui um Work Package a você.

```bash
op wp assign-me 123
```

### `op wp create-from-image`

Cria um Work Package a partir de uma imagem usando IA local (Ollama).

```bash
op wp create-from-image ./screenshot.png
op wp create-from-image ./erro.jpg --model llava
op wp create-from-image ./bug.png -y
op wp create-from-image --clipboard
```

| Flag | Alias | Descrição |
|------|-------|-----------|
| `--model` | `-m` | modelo ollama para análise (default: llava) |
| `--yes` | `-y` | criar sem pedir confirmação |
| `--clipboard` | `-c` | usar imagem do clipboard |

**Requisitos:** Ollama rodando localmente com um modelo de visão (llava, minicpm-v, etc.)

```bash
ollama serve
ollama pull llava
```

## Roadmap

Ideias em desenvolvimento:

- mais comandos para gerenciamento de projetos
- filtros avançados na listagem
- integração com git para criar tasks automaticamente

## Contribuindo

Toda ideia de melhoria é bem-vinda! Sinta-se à vontade para abrir issues ou enviar PRs.

## Autor

Desenvolvido por [Guilherme Alves](https://www.guialves.site)
