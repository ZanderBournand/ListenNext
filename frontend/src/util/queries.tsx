import { gql } from "@apollo/client";

export const SpotifyLoginUrl = `
  query SpotifyLoginUrl {
    spotifyUrl
  }
`

export const LastScrapeTime = `
  query LastScrapeTime {
    lastScrapeTime
  }
`

export const querySearchArtists = `
  query SearchArtists($query: String!){
    searchArtists(query: $query) {
      results{
        spotify_id
        name
        image
        popularity
        recent_releases_count
        upcoming_releases_count
      }
    	related_artists{
      	spotify_id
        name
        image
        popularity
        recent_releases_count
        upcoming_releases_count
    	}
    }
  }
`

export const queryAllRecommendations = `
  query AllRecommendations {
    allRecommendations {
      artists {
        name
        popularity
        spotify_id
      }
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
`

export const queryReleasesCount = gql`
  query GetAllReleasesCount{
    allReleasesCount{
      past{
        all
        albums
        singles
      }
      week{
        all
        albums
        singles
      }
      month{
        all
        albums
        singles
      }
      extended{
        all
        albums
        singles
      }
    }
  }
`;

export const queryTrendingReleases = gql`
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
      next
    }
  }
`;

export const queryAllTrendingReleases = ({ releaseType }: any) => gql`
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