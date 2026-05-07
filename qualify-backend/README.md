# Sobre
------- 
Este é o *backend* do  Qualify.

# Primeiros Passos
---------

1. Instale o Docker e Docker Compose e adicione seu usuário ao grupo de usuários que ele cria. Em Arch Linux, isso pode ser feito assim:
```
sudo pacman -S docker docker-compose
sudo usermod -aG $USER docker
```

2. Entre no diretório do projeto:
```
cd qualify-backend
```

3. Dê build no docker:
```
docker build --tag qualify-backend .
```

4. Rode o docker (ele foi configurado na porta `3001`):
```
docker run -p 3001:3001 qualify-backend
```

Alterações no projeto requerem que o comando do passo `2` seja executado novamente após as mudanças. O projeto vai estar rodando em:
```
http://localhost:3001/
```