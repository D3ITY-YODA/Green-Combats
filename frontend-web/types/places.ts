import { UUID } from "./common";
export type PlaceType = "country" | "region" | "district" | "ward" | "village" | "community" | "watershed" | "waterbody" | "custom_area";
export interface Place {
  id: UUID;
  name: string;
  type: PlaceType;
  country_code: string;
}
export interface UserPlace {
  id: UUID;
  user_id: UUID;
  place_id: UUID;
  label?: string;
  is_primary: boolean;
  place: Place;
}
