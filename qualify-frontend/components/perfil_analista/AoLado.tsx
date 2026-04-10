interface perfilCardsInter {
    dado: string;
    descricao: string;
}

const perfilTaxas: perfilCardsInter[] = [
    {
        dado: "124",
        descricao: "Projetos Completos"
    },
    {
        dado: "98%",
        descricao: "Taxa de sucesso"
    },
    {
        dado: "R$120/hr",
        descricao: "Preço por hora"
    },
    {
        dado: "1k+",
        descricao: "Bugs descobridos"
    },
];

function PerfilCards({ dado, descricao }: perfilCardsInter) {
  return (
    <div className="flex justify-start inset-ring-2 inset-ring-zinc-800 flex-col w-full mx-3 mt-4 bg-gray-950">
        <p className="text-blue-800 text-2xl mt-4 mx-4">
            {dado}
        </p>
        <p className="text-zinc-400 mt-1 mb-4 mx-4">
            {descricao}
        </p>
    </div>
  );
}

export function AoLado(){
    return(
        <div className="flex flex-col justify-between content-center mt-8 w-2/10 h-100">
            {perfilTaxas.map((perfilm) => (
                <PerfilCards key={perfilm.descricao} {...perfilm} />
            ))}
        </div>
    );
}