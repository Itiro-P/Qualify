"use client";

import { Alert } from "@/components/ui/Alert";
import { Loading } from "@/components/ui/Loading";
import { analystService } from "@/libs/services/analystService";
import { Analyst } from "@/types/services/analyst";
import { Review } from "@/types/services/review";
import { Stars } from "lucide-react";
import { useEffect, useState } from "react";

function StarRange({ rating }: { rating: number }) {
  const StarRange = [];

  for (let i = 1; i <= 5; i++) {
    if (rating >= i) {
      StarRange.push(<Stars key={i} type="full" />);
    } else {
      StarRange.push(<Stars key={i} type="empty" />);
    }
  }

  return <div>{StarRange}</div>;
}

function ListServicesArray({ reviews }: { reviews: Review[] }) {
  return (
    <div>
      {reviews.map((review) => (
        <div key={review.id} className="flex flex-col p-4 mb-6">
          <p className="text-zinc-200">{review.comment}</p>
          <div className="flex text-zinc-400">
            <div className="my-4 mr-4">
              <StarRange rating={review.rating} />
            </div>
            <p className="my-4 ml-4">{review.time_created}</p>
          </div>
        </div>
      ))}
    </div>
  );
}

export function ListReviews({ analyst }: { analyst: Analyst }) {
  const [reviewsAnalyst, setReviewsAnalyst] = useState<Review[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    async function getInfo() {
      setLoading(true);
      const analystReviews = await analystService.listReviews(analyst.id);
      if (analystReviews) {
        setReviewsAnalyst(analystReviews);
      } else {
        setError("Erro ao carregar as avaliações do analista.");
      }
      setLoading(false);
    }
    getInfo();
  }, [analyst.id]);

  return loading ? (
    <Loading />
  ) : (
    <div>
      <div className="flex">
        {error && <Alert variant="error">{error}</Alert>}
        <ListServicesArray reviews={reviewsAnalyst} />
      </div>
    </div>
  );
}
