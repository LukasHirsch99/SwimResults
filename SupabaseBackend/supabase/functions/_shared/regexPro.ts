export const evaluteRegex = (r: RegExp, s: string) =>
{
  const matches = s.matchAll(r)
  const matchDict: Map<string, string> = new Map()

  for (const match of matches)
  {
    if (match.groups)
      Object.entries(match.groups).forEach(([k, v]) =>
      {
        if (v)
          matchDict.set(k, v)
      })
  }
  return matchDict
}