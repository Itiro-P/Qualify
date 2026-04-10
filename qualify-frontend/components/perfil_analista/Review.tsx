import { SquarePen } from 'lucide-react';

interface reviewsCardsAnal {
    avaliacao: number,
    comentario: string,
    data: string,
}

const reviewsAnal: reviewsCardsAnal[] = [
    {
        avaliacao: 2.5,
        comentario: "Bons testes aletórios",
        data: "10/03/2026",
    },
    {
        avaliacao: 3.2,
        comentario: "É um exelente guardião da qualidade",
        data: "21/01/2026",
    },
    {
        avaliacao: 2.7,
        comentario: "Auxiliou muito com seus testes de botões",
        data: "15/02/2026",
    },
];

function Estrelas({ type }: {type: string}) {
    if (type === "full") return <span>⭐</span>;
    if (type === "half") return <span>🌓</span>;
    return <span>☆</span>;
};

function RangeEstrelas ({ rating }: {rating: number}) {
    const EstrelasRange = [];

    for (let i = 1; i <= 5; i++) {
        if (rating >= i) {
            EstrelasRange.push(<Estrelas key={i} type="full" />);
        } else if (rating >= i - 0.5) {
            EstrelasRange.push(<Estrelas key={i} type="half" />);
    } else {
        EstrelasRange.push(<Estrelas key={i} type="empty" />);
    }
}

return <div>{EstrelasRange}</div>;
};

function ReviewCards({ avaliacao, comentario, data }: reviewsCardsAnal ) {
    return (
        <div id="comentario" className="flex flex-col p-4 mb-6">
            <p className="text-zinc-200">{comentario}</p>
            <div className="flex text-zinc-400">
                <div className="my-4 mr-4">
                    <RangeEstrelas rating={avaliacao}/>
                </div>
                <p className="my-4 ml-4">{data}</p>
            </div>
        </div>
    );
}

export function Review(){
    return (
        <div id="review" className="flex flex-col p-4 mb-6">
            <div className="flex">
                <SquarePen className="text-blue-500 mr-2"/>
                <h2 className="text-xl">
                    Reviews
                </h2>
            </div>
            <div className="px-2 py-1 m-1">
                {reviewsAnal.map((reviews) => (
                    <ReviewCards key={reviews.comentario} {...reviews} />
                ))}
            </div>
        </div>
    );
}
