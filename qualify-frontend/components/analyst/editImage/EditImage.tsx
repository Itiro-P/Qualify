"use client";

import { useEffect, useState } from "react";
import { FormInput } from "@/components/ui";
import { Analyst } from "@/types/services";
import { ImageProfile } from "@/components/analyst/profile";
import { analystService } from "@/libs";

function handleChange(
  e: React.ChangeEvent<HTMLInputElement>,
  setForm: React.Dispatch<React.SetStateAction<File | undefined>>,
) {
  setForm(e.target.files?.[0]);
}

export function EditImage({
  analyst,
  setAnalystImage,
}: {
  analyst: Analyst;
  setAnalystImage: React.Dispatch<React.SetStateAction<File | undefined>>;
}) {
  const [imageUrl, setImageUrl] = useState<string>();

  useEffect(() => {
    async function loadImage() {
      const response = await analystService.getProfile(analyst.id);
      setImageUrl(response?.picture);
    }

    loadImage();
  }, [analyst.id]);

  return (
    <div>
      <ImageProfile user={analyst.name} imageURL={imageUrl} />

      <p className="mb-2 text-sm font-medium text-white/80">Foto de perfil</p>

      <form className="mt-2 mb-8 rounded-xl border border-white/10 p-5 pt-3">
        <div className="flex flex-col gap-4">
          <FormInput
            label="Imagem"
            name="image"
            type="file"
            accept="image/*"
            onChange={(e) => handleChange(e, setAnalystImage)}
            required
          />
        </div>
      </form>
    </div>
  );
}
