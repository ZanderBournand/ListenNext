import { gql } from "@apollo/client";

export const query = gql`
  query GetTrendingReleases($releaseType: String!, $direction: String!, $reference: Int!, $period: String!) {
    trendingReleases(input: {type: $releaseType, direction: $direction, reference: $reference, period: $period}) {
      releases {
        _id
        title
        artists {
          name
        }
        release_date
        cover
        type
        trending_score
      }
    }
  }
`;

export const queryAll = ({ releaseType }: any) => gql`
  query GetAllTrendingReleases {
    allTrendingReleases(type: "${releaseType}") {
        past {
            _id
            title
            artists{
            name
            }
            release_date
            cover
            type
            trending_score
        }
        week {
            _id
            title
            artists{
            name
            }
            release_date
            cover
            type
            trending_score
            aoty_id
        }
        month {
            _id
            title
            artists{
            name
            }
            release_date
            cover
            type
            trending_score
        }
        extended {
            _id
            title
            artists{
            name
            }
            release_date
            cover
            type
            trending_score
        }
    }
  }
`;