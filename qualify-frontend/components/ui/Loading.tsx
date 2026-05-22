export function Loading() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-950">
      <div className="flex flex-col items-center gap-4">
        <div className="relative h-16 w-16">
          <div className="absolute inset-0 animate-ping rounded-full bg-blue-500 opacity-30"></div>
          <div className="absolute inset-2 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"></div>
        </div>

        <div className="space-y-2 text-center">
          <h1 className="text-xl font-semibold text-white">Carregando...</h1>
          <p className="text-sm text-gray-400">
            Aguarde enquanto buscamos as informações.
          </p>
        </div>
      </div>
    </div>
  );
}
