import { Settings } from 'lucide-react';

interface tecnologiaCardsAnal {
    tecnologia: string;
}

const tecnologiasAnal: tecnologiaCardsAnal[] = [
    {
        tecnologia: "Docker",
    },
    {
        tecnologia: "Javascript/TS",
    },
    {
        tecnologia: "GitHub Actions",
    },
    {
        tecnologia: "Kubernets",
    },
];

function TecnologiasCards({ tecnologia }: tecnologiaCardsAnal) {
    return (
        <p className="text-blue-600 border-blue-600 rounded-xl text-lg px-2 py-1 m-1 bg-blue-950">
            {tecnologia}
        </p>
    );
}

export function Tecnologias(){
    return (
        <div className="flex flex-col p-4 mb-6">
            <div className="flex">
                <Settings className="text-blue-500 mr-2"/>
                <h2 className="text-xl">
                    Tecnologias
                </h2>
            </div>
            <div className="flex mt-2 mx-2 ">
                {tecnologiasAnal.map((tecnologias) => (
                    <TecnologiasCards key={tecnologias.tecnologia} {...tecnologias} />
                ))}
            </div>
        </div>
    );
}