#!/usr/bin/env python3
import argparse
import requests
import json
import sys
import os
from typing import Dict, Any

TAVILY_API_KEY = os.getenv("TAVILY_API_KEY")
if not TAVILY_API_KEY:
    print("Error: TAVILY_API_KEY environment variable is not set", file=sys.stderr)
    sys.exit(1)

API_URL = "https://api.tavily.com/search"

def search(query: str, search_depth: str = "basic") -> Dict[Any, Any]:
    headers = {"api-key": TAVILY_API_KEY}
    params = {
        "query": query,
        "search_depth": search_depth,
        "include_domains": [],
        "exclude_domains": []
    }
    
    try:
        response = requests.post(API_URL, json=params, headers=headers)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)

def format_results(data: Dict[Any, Any]) -> None:
    for result in data['results']:
        print(f"\n🔗 {result['url']}")
        print(f"📝 {result['content']}\n")
        print("-" * 80)

def main():
    parser = argparse.ArgumentParser(description="Tavily Search CLI")
    parser.add_argument("query", help="Search query")
    parser.add_argument("--json", action="store_true", help="Output in JSON format")
    parser.add_argument("--depth", choices=["basic", "advanced"], default="basic", 
                      help="Search depth (basic or advanced)")
    
    args = parser.parse_args()
    
    results = search(args.query, args.depth)
    
    if args.json:
        print(json.dumps(results, indent=2))
    else:
        format_results(results)

if __name__ == "__main__":
    main() 