import { api } from "@/libs/api";
import type {
  Review,
  ReviewResponse,
  ReviewsResponse,
  ReviewUpdateRequest,
  ListReviewsParams,
} from "@/types/services/review";

function buildQuery(params?: Record<string, unknown>): string {
  if (!params) return "";
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== null) {
      search.set(key, String(value));
    }
  }
  const qs = search.toString();
  return qs ? `?${qs}` : "";
}

export const reviewService = {
  list(params?: ListReviewsParams): Promise<ReviewsResponse> {
    return api.get(`/reviews${buildQuery(params as Record<string, unknown>)}`);
  },

  getById(id: number): Promise<Review | null> {
    return api.get<ReviewResponse>(`/reviews/${id}`).then(
      (resp) => {
        return resp.review;
      },
      () => {
        return null;
      },
    );
  },

  create(data: Review): Promise<Review | null> {
    return api.post<ReviewResponse>("/reviews", data).then(
      (resp) => {
        return resp.review;
      },
      () => {
        return null;
      },
    );
  },

  update(id: number, data: Review): Promise<Review | null> {
    return api.put<ReviewResponse>(`/reviews/${id}`, data).then(
      (resp) => {
        return resp.review;
      },
      () => {
        return null;
      },
    );
  },

  patch(id: number, data: ReviewUpdateRequest): Promise<Review | null> {
    return api.patch<ReviewResponse>(`/reviews/${id}`, data).then(
      (resp) => {
        return resp.review;
      },
      () => {
        return null;
      },
    );
  },

  delete(id: number): Promise<Record<string, string>> {
    return api.delete(`/reviews/${id}`);
  },
};
