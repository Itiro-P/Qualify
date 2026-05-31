import { FormInput } from "@/components/ui";
import { IFormResponse } from "@/types/customer/formResponse";
import { Dispatch, SetStateAction } from "react";

export function SearchForm({
  formResponse,
  setFormResponse,
}: {
  formResponse: IFormResponse | null;
  setFormResponse: Dispatch<SetStateAction<IFormResponse | null>>;
}) {
  return (
    <form>
      <h3>Habilidades</h3>
      <FormInput
        label="Habilidades"
        onChange={(e) =>
          setFormResponse({
            value: formResponse?.value || 0,
            skills: formResponse
              ? [...formResponse.skills, e.target.value]
              : [e.target.value],
            localization: formResponse?.localization || {
              country: "",
              state: "",
              city: "",
            },
          })
        }
      />
      <h3>Preço por hora</h3>
      <FormInput
        label="Preço"
        type="range"
        min="0"
        max="5000"
        value={formResponse?.value || 0}
        onChange={(e) =>
          setFormResponse({
            value: Number(e.target.value),
            skills: formResponse?.skills || [],
            localization: formResponse?.localization || {
              country: "",
              state: "",
              city: "",
            },
          })
        }
      />
      <h3>Localização</h3>
      <FormInput
        label="País"
        value={formResponse?.localization?.country || ""}
        onChange={(e) =>
          setFormResponse({
            value: formResponse?.value || 0,
            skills: formResponse?.skills || [],
            localization: formResponse?.localization
              ? {
                  ...formResponse.localization,
                  country: e.target.value,
                }
              : {
                  country: e.target.value,
                  state: "",
                  city: "",
                },
          })
        }
      />
      <FormInput
        label="Estado"
        value={formResponse?.localization?.state || ""}
        onChange={(e) =>
          setFormResponse({
            value: formResponse?.value || 0,
            skills: formResponse?.skills || [],
            localization: formResponse?.localization
              ? {
                  ...formResponse.localization,
                  state: e.target.value,
                }
              : {
                  country: "",
                  state: e.target.value,
                  city: "",
                },
          })
        }
      />
      <FormInput
        label="Cidade"
        value={formResponse?.localization?.city || ""}
        onChange={(e) =>
          setFormResponse({
            value: formResponse?.value || 0,
            skills: formResponse?.skills || [],
            localization: formResponse?.localization
              ? {
                  ...formResponse.localization,
                  city: e.target.value,
                }
              : {
                  country: "",
                  state: "",
                  city: e.target.value,
                },
          })
        }
      />
    </form>
  );
}
