const std = @import("std");

const Store = @import("../core/store.zig").Store;
const command = @import("../core/command.zig");
const engine = @import("../core/engine.zig");
const response = @import("../core/response.zig");
const clock = @import("../core/clock.zig");

pub const App = struct {
    allocator: std.mem.Allocator,
    store: Store,

    pub fn init(allocator: std.mem.Allocator) App {
        return .{
            .allocator = allocator,
            .store = Store.init(allocator),
        };
    }

    pub fn deinit(self: *App) void {
        self.store.deinit();
    }

    pub fn executeText(self: *App, input: []const u8) ![]u8 {
        const cmd = command.parse(input) catch |err| {
            return response.err(self.allocator, @errorName(err));
        };

        var exec = engine.Engine.init(&self.store);
        const fixed_clock = clock.Clock.fixed(0);

        return exec.executeAt(self.allocator, cmd, fixed_clock.now());
    }
};

test "app executes ping" {
    var app = App.init(std.testing.allocator);
    defer app.deinit();

    const resp = try app.executeText("PING");
    defer std.testing.allocator.free(resp);

    try std.testing.expectEqualStrings("+PONG\r\n", resp);
}

test "app returns parse error response" {
    var app = App.init(std.testing.allocator);
    defer app.deinit();

    const resp = try app.executeText("UNKNOWN");
    defer std.testing.allocator.free(resp);

    try std.testing.expectEqualStrings("-ERR UnknownCommand\r\n", resp);
}

test "app executes set and get in same instance" {
    var app = App.init(std.testing.allocator);
    defer app.deinit();

    {
        const resp = try app.executeText("SET name zigkv");
        defer std.testing.allocator.free(resp);
        try std.testing.expectEqualStrings("+OK\r\n", resp);
    }

    {
        const resp = try app.executeText("GET name");
        defer std.testing.allocator.free(resp);
        try std.testing.expectEqualStrings("$zigkv\r\n", resp);
    }
}
