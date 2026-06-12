"use client";

import { Alert } from "@/components/ui/Alert";
import { Loading } from "@/components/ui/Loading";
import { reviewService } from "@/libs/services/reviewService";
import { Review } from "@/types/services/review";
import { User } from "@/types/services/user";
import { Star, MessageSquare } from "lucide-react";
import { useEffect, useState } from "react";

function StarRange({ rating }: { rating: number }) {
  return (
    <div className="flex items-center gap-0.5">
      {Array.from({ length: 5 }, (_, i) => (
        <Star
          key={i}
          className={`w-4 h-4 ${
            rating >= i + 1
              ? "text-accent fill-accent"
              : "text-white/20 fill-transparent"
          }`}
        />
      ))}
    </div>
  );
}

function formatDate(value?: string) {
  if (!value) return "";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleDateString("pt-BR", {
    day: "2-digit",
    month: "long",
    year: "numeric",
  });
}

function ReviewCard({ review }: { review: Review }) {
  return (
    <div className="rounded-xl border border-white/10 bg-white/5 p-6">
      <div className="flex items-center justify-between gap-4">
        <StarRange rating={review.rating} />
        <span className="text-xs text-neutral-slate">
          {formatDate(review.time_created)}
        </span>
      </div>
      <p className="mt-3 text-sm leading-relaxed text-slate-200">
        {review.comment || "Sem comentário."}
      </p>
    </div>
  );
}

export function ListReviews({ user }: { user: User }) {
  const [reviewsAnalyst, setReviewsAnalyst] = useState<Review[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    async function getInfo() {
      setLoading(true);
      const reviews = await reviewService.listReviewsByClient(user.id);
      if (reviews) {
        setReviewsAnalyst(reviews);
      } else {
        setError("Erro ao carregar as avaliações do usuário.");
      }
      setLoading(false);
    }
    getInfo();
  }, [user.id]);

  if (loading) return <Loading />;

  const average =
    reviewsAnalyst.length > 0
      ? reviewsAnalyst.reduce((sum, r) => sum + r.rating, 0) /
        reviewsAnalyst.length
      : 0;

  return (
    <section className="px-6 md:px-20 py-14">
      <div className="max-w-4xl mx-auto">
        <div className="mb-10">
          <h1 className="text-3xl font-bold mb-4">Avaliações</h1>
          <div className="h-1 w-20 bg-accent mb-6" />
          {reviewsAnalyst.length > 0 && (
            <div className="flex items-center gap-3">
              <StarRange rating={Math.round(average)} />
              <span className="text-lg font-semibold text-white">
                {average.toFixed(1)}
              </span>
              <span className="text-sm text-neutral-slate">
                ({reviewsAnalyst.length}{" "}
                {reviewsAnalyst.length === 1 ? "avaliação" : "avaliações"})
              </span>
            </div>
          )}
        </div>

        {error && <Alert variant="error">{error}</Alert>}

        {!error && reviewsAnalyst.length === 0 ? (
          <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-white/10 bg-white/2 py-20 text-center">
            <MessageSquare className="w-10 h-10 text-neutral-slate mb-4" />
            <p className="text-neutral-slate">
              Este analista ainda não possui avaliações.
            </p>
          </div>
        ) : (
          <div className="flex flex-col gap-4">
            {reviewsAnalyst.map((review) => (
              <ReviewCard key={review.id} review={review} />
            ))}
          </div>
        )}
      </div>
    </section>
  );
}
