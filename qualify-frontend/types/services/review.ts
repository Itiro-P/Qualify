export interface Review {
  id?: number;
  analyst_id?: number;
  client_id?: number;
  service_id?: number;
  rating?: number;
  comment?: string;
  time_created?: string;
}

export interface ReviewResponse {
  review: Review;
}

export interface ReviewsResponse {
  reviews: Review[];
  count: number;
  page?: number;
  page_size?: number;
}

export interface ReviewUpdateRequest {
  rating?: number;
  comment?: string;
}

export interface ListReviewsParams {
  analyst_id?: number;
  client_id?: number;
  service_id?: number;
  rating?: number;
  min_rating?: number;
  max_rating?: number;
  page?: number;
  page_size?: number;
}
