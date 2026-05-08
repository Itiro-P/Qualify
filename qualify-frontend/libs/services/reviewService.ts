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

  getById(id: number): Promise<ReviewResponse> {
    return api.get(`/reviews/${id}`);
  },

  create(data: Review): Promise<ReviewResponse> {
    return api.post("/reviews", data);
  },

  update(id: number, data: Review): Promise<ReviewResponse> {
    return api.put(`/reviews/${id}`, data);
  },

  patch(id: number, data: ReviewUpdateRequest): Promise<ReviewResponse> {
    return api.patch(`/reviews/${id}`, data);
  },

  delete(id: number): Promise<Record<string, string>> {
    return api.delete(`/reviews/${id}`);
  },
};
