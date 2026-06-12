"use client";

import { SearchForm, SearchResult } from "@/components/customer/searchAnalyst";
import { IFormResponse } from "@/types/customer/formResponse";
import { useState } from "react";

export function SearchAnalyst() {
  const [formResponse, setFormResponse] = useState<IFormResponse | null>(null);

  return (
    <div>
      <div className="p-4 border rounded-md mb-4 max-w-2/10">
        <SearchForm
          formResponse={formResponse}
          setFormResponse={setFormResponse}
        />
      </div>
      <div>
        <SearchResult formResponse={formResponse} />
      </div>
    </div>
  );
}
