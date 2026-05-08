import { ITechnology } from "@/types/analyst/profile/technology";

export interface IRegisterTechnology {
  technologiesAnalyst: ITechnology[];
  setTechnologiesAnalyst: React.Dispatch<React.SetStateAction<ITechnology[]>>;
}
