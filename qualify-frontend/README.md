# Sobre
---------
Este é o *frontend* do Qualify.

# Primeiros Passos
---------

1. Instale as dependências:
```bash
npm install
```

2. Rode o projeto:
- Para a *build* de desenvolvimento: `npm run dev`;
- Para a *build* de produção: `npm run build` e `npm run start`;

Enquanto `npm run dev` estiver em execução. Alterar qualquer arquivo resultará em mudanças imediatas na tela.


# Arquitetura do projeto
--------
- Em `app/` os arquivos `.tsx` ficam organizados. Importe com `@app`;
- Em `styles` os estilos *css*. Importe com `@styles`;
- Em `public` ficam quaisquer *assets* que sejam necessários. Importe com `@public`;
