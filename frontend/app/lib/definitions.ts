// Type for what gets pulled from the database.
export type VendorContract = {
  seller: string;
  effective_date: string;
  renewal_date: string;
  autorenew: boolean;
};
