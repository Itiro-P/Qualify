export interface ProposalLetter {
  id?: number;
  analyst_id?: number;
  client_id?: number;
  title?: string;
  content?: string;
  proposed_hourly_rate?: number;
  time_created?: string;
}

export interface ProposalLetterResponse {
  proposal_letter: ProposalLetter;
}

export interface ProposalLettersResponse {
  proposal_letters: ProposalLetter[];
  count: number;
}

export interface ProposalLetterUpdateRequest {
  title?: string;
  content?: string;
  proposed_hourly_rate?: number;
}

export interface ListProposalsParams {
  client_id?: number;
  analyst_id?: number;
}
