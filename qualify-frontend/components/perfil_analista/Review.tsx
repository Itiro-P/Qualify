import { SquarePen } from 'lucide-react';

interface reviewsCardsInterface {
    rating: number,
    comment: string,
    date: string,
}

const reviewCardsVector: reviewsCardsInterface[] = [
    {
        rating: 2.5,
        comment: "Bons testes aletórios",
        date: "10/03/2026",
    },
    {
        rating: 3.2,
        comment: "É um exelente guardião da qualidade",
        date: "21/01/2026",
    },
    {
        rating: 2.7,
        comment: "Auxiliou muito com seus testes de botões",
        date: "15/02/2026",
    },
];

function Stars({ type }: {type: string}) {
    if (type === "full") return <span>⭐</span>;
    return <span>☆</span>;
};

function StarRange ({ rating }: {rating: number}) {
    const StarRange = [];

    for (let i = 1; i <= 5; i++) {
        if (rating >= i) {
            StarRange.push(<Stars key={i} type="full" />);
        } else {
            StarRange.push(<Stars key={i} type="empty" />);
        }
    }

    return <div>{StarRange}</div>;
};

function ReviewCards({ rating, comment, date }: reviewsCardsInterface ) {
    return (
        <div id="comment" className="flex flex-col p-4 mb-6">
            <p className="text-zinc-200">{comment}</p>
            <div className="flex text-zinc-400">
                <div className="my-4 mr-4">
                    <StarRange rating={rating}/>
                </div>
                <p className="my-4 ml-4">{date}</p>
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
                {reviewCardsVector.map((reviews) => (
                    <ReviewCards key={reviews.comment} {...reviews} />
                ))}
            </div>
        </div>
    );
}
