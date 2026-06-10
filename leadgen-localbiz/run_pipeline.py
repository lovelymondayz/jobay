
"""
Master Pipeline Runner — runs all 7 stages in sequence.
Usage: python3 run_pipeline.py [--stage scout|evaluate|enrich|build|deploy|outreach|track]
"""
import sys, importlib, time

STAGES = {
    "scout": "scripts.scout.find_businesses",
    "evaluate": "scripts.evaluate.check_websites",
    "enrich": "scripts.enrich.enrich_data",
    "build": "scripts.build.generate_page",
    "deploy": "scripts.deploy.push_deploy",
    "outreach": "scripts.outreach.send_pitches",
    "track": "scripts.track.daily_report",
}

def run_stage(stage_name):
    module_path = STAGES[stage_name]
    print(f"\n{'='*60}")
    print(f"Running stage: {stage_name}")
    print(f"{'='*60}")
    mod = importlib.import_module(module_path)
    mod.main()

def main():
    if len(sys.argv) > 1 and sys.argv[1] == "--stage":
        stage = sys.argv[2]
        if stage not in STAGES:
            print(f"Unknown stage: {stage}")
            print(f"Available: {', '.join(STAGES.keys())}")
            sys.exit(1)
        run_stage(stage)
    else:
        # Run all stages in order
        start = time.time()
        for stage_name in STAGES:
            try:
                run_stage(stage_name)
            except Exception as e:
                print(f"ERROR in {stage_name}: {e}")
                # Continue to next stage
        elapsed = time.time() - start
        print(f"\n{'='*60}")
        print(f"Pipeline complete in {elapsed:.1f}s")
        print(f"{'='*60}")

if __name__ == "__main__":
    main()
