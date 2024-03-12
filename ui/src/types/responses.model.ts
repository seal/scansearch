export interface ISearchResponse {
  ShoppingResults: IShoppingResult[] | null;
  filters: ISearchFilter[] | null;
  message: string;
  success: boolean;
}
export interface IShoppingResult {
  delivery: string;
  extensions: string[];
  extracted_price: number;
  link: string;
  position: number;
  price: string;
  product_id: string;
  product_link: string;
  rating: number;
  reviews: number;
  serpapi_product_api: string;
  source: string;
  thumbnail: string;
  title: string;
  isActive: boolean;
}

export interface ISearchFilter {
  options: {
    tbs: string;
    text: string;
  }[];
  type: string;
}

export interface IUserDetails {
  ID: number;
  firstname: string;
  lastname: string;
  username: string;
  email: string;
  plan: "free";
}

export interface ILoginResponse {
  token: string;
  user: IUserDetails;
}

export interface IWardrobeItem {
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt?: string;
  ID: number;
  userid: number;
  serpapi_product_api: string;
  position: number;
  title: string;
  link: string;
  product_link: string;
  source: string;
  price: string;
  extracted_price: number;
  rating: number;
  thumbnail: string;
  delivery: string;
}
