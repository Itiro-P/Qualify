"use client";

import { SearchForm, SearchResult } from "@/components/customer/searchAnalyst";
import { IFormResponse } from "@/types/customer/formResponse";
import { useState } from "react";

export function SearchAnalyst() {
  const [formResponse, setFormResponse] = useState<IFormResponse | null>(null);

  return (
    <div>
      <SearchForm
        formResponse={formResponse}
        setFormResponse={setFormResponse}
      />
      <SearchResult
        formResponse={formResponse}
        setFormResponse={setFormResponse}
      />
    </div>
  );
}
