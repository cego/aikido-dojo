# Excluding SAST Findings Using Comments

It is possible to ignore specific SAST finding using a `//nosec` comment at the first line or preceding line of the pattern match. By adding this comment, the finding will be ignored in the Aikido feed and marked as *manually suppressed by developer.*

Below you will find some examples of how a NoSQL injection issue can be ignored using comments.

Ignore by adding comment to line preceding the detected issue:

```typescript
module.exports = function trackOrder () {
    return (req: Request, res: Response) => {
        const id = foo.a() ? String(req.params.id) : req.params.id
        // nosec
        Solution.findOne({ _id: req.body.id }) 
    }
}
```

Ignore by adding comment to first (or only) of the detected issue:

```typescript
module.exports = function trackOrder () {
    return (req: Request, res: Response) => {
        const id = foo.a() ? String(req.params.id) : req.params.id
        Solution.findOne({ _id: req.body.id }) // nosec
    }
}
```
